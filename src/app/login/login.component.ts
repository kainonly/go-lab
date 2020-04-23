import {Component, OnInit} from '@angular/core';
import {Router} from '@angular/router';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';
import {NzNotificationService} from 'ng-zorro-antd';
import {MainService} from '@common/main.service';
import {BitService, StorageService} from 'ngx-bit';
import {StorageMap} from '@ngx-pwa/local-storage';
import {switchMap} from 'rxjs/operators';
import packer from './language';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {
  form: FormGroup;
  users: any[] = [];

  constructor(
    private mainService: MainService,
    private notification: NzNotificationService,
    private router: Router,
    private fb: FormBuilder,
    public bit: BitService,
    private storageService: StorageService,
    private storageMap: StorageMap
  ) {
  }

  ngOnInit() {
    this.bit.registerLocales(packer);
    this.form = this.fb.group({
      username: [null, [
        Validators.required,
        Validators.minLength(4),
        Validators.maxLength(20)
      ]],
      password: [null, [
        Validators.required,
        Validators.minLength(8),
        Validators.maxLength(18)]
      ],
      remember: [1, [
        Validators.required
      ]]
    });
    this.storageMap.get('users').subscribe((data: Set<string>) => {
      if (data) {
        this.users = [...data.keys()];
      }
    });
  }

  submit(data: any) {
    this.mainService.login(data.username, data.password).subscribe(res => {
      if (!res.error) {
        this.storageService.clear();
        if (data.remember) {
          this.storageMap.get('users').pipe(
            switchMap((lists: Set<string>) =>
              this.storageMap.set(
                'users',
                lists ? lists.add(data.username) : new Set([data.username])
              )
            )
          ).subscribe(() => {
          });
        }
        sessionStorage.setItem('currentUsername', data.username);
        sessionStorage.setItem('loginTime', new Date().toISOString());
        this.notification.success(
          this.bit.l.loginTips,
          this.bit.l.loginSuccess
        );
        this.router.navigateByUrl('/');
      } else {
        this.notification.error(
          this.bit.l.loginTips,
          this.bit.l.loginFailed
        );
      }
    });
  }
}