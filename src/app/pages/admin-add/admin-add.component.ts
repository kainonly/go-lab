import {Component, OnInit} from '@angular/core';
import {AbstractControl, FormBuilder, FormGroup, Validators} from '@angular/forms';
import {SwalService, BitService, asyncValidator} from 'ngx-bit';
import {switchMap} from 'rxjs/operators';
import {NzNotificationService} from 'ng-zorro-antd';
import {AdminService} from '@common/admin.service';
import {RoleService} from '@common/role.service';
import packer from './language';

@Component({
  selector: 'app-admin-add',
  templateUrl: './admin-add.component.html',
})
export class AdminAddComponent implements OnInit {
  form: FormGroup;
  avatar = '';
  roleLists: any[] = [];

  constructor(
    private swal: SwalService,
    private fb: FormBuilder,
    public bit: BitService,
    private notification: NzNotificationService,
    private adminService: AdminService,
    private roleService: RoleService
  ) {
  }

  ngOnInit() {
    this.bit.registerLocales(packer);
    this.form = this.fb.group({
      username: [null, [
        Validators.required,
        Validators.minLength(4),
        Validators.maxLength(20)
      ], [this.validedUsername]
      ],
      password: [null, [
        Validators.required,
        Validators.minLength(8),
        Validators.maxLength(18)]
      ],
      role: [null, [Validators.required]],
      call: [null],
      email: [null, [Validators.email]],
      phone: [null],
      status: [true, [Validators.required]]
    });
    this.getRole();
  }

  validedUsername = (control: AbstractControl) => {
    return asyncValidator(this.adminService.validedUsername(control.value));
  };

  /**
   * 获取权限组
   */
  getRole() {
    this.roleService.originLists().subscribe(data => {
      this.roleLists = data;
    });
  }

  /**
   * 上传
   */
  upload(info) {
    if (info.type === 'success') {
      this.avatar = info.file.response.data.save_name;
      this.notification.success(
        this.bit.l.success,
        this.bit.l.uploadSuccess
      );
    }
    if (info.type === 'error') {
      this.notification.error(
        this.bit.l.notice,
        this.bit.l.uploadError
      );
    }
  }

  /**
   * 提交
   */
  submit(data) {
    if (this.avatar) {
      data.avatar = this.avatar;
    }
    this.adminService.add(data).pipe(
      switchMap(res => this.swal.addAlert(res, this.form, {
        status: true
      }))
    ).subscribe(() => {
    });
  }
}
