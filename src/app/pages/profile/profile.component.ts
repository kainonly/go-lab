import {Component, OnInit} from '@angular/core';
import {BitService} from 'ngx-bit';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';
import {NzNotificationService} from 'ng-zorro-antd';
import {MainService} from '@common/main.service';
import packer from './language';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
})
export class ProfileComponent implements OnInit {
  form: FormGroup;
  editPassword = false;
  avatar = '';

  constructor(
    public bit: BitService,
    private fb: FormBuilder,
    private mainService: MainService,
    private notification: NzNotificationService
  ) {
  }

  ngOnInit() {
    this.bit.registerLocales(packer);
    this.form = this.fb.group({
      call: [null],
      email: [null, [Validators.email]],
      phone: [null],
      old_password: [null, [Validators.minLength(8), Validators.maxLength(18)]],
      new_password: [null, [Validators.minLength(8), Validators.maxLength(18)]]
    });
    this.getInformation();
  }

  /**
   * 获取个人信息
   */
  getInformation() {
    this.mainService.information().subscribe(data => {
      this.avatar = data.avatar;
      this.form.setValue({
        call: data.call,
        email: data.email,
        phone: data.phone,
        old_password: null,
        new_password: null
      });
    });
  }

  /**
   * 头像上传
   */
  upload(info) {
    if (info.type === 'success') {
      this.avatar = info.file.response.data.savename;
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
   * 监听密码修改关闭
   */
  editPasswordChange(status: boolean) {
    if (!status) {
      this.form.get('old_password')
        .setValue(null);
      this.form.get('new_password')
        .setValue(null);
    }
  }

  /**
   * 提交
   */
  submit(data) {
    if (this.avatar) {
      data.avatar = this.avatar;
    }
    if (!this.editPassword) {
      delete data.old_password;
      delete data.new_password;
    }
    this.mainService.update(data).subscribe(res => {
      if (res.error) {
        if (res.msg === 'error:password') {
          this.notification.error(
            this.bit.l.failed,
            this.bit.l.passwordError
          );
        } else {
          this.notification.error(
            this.bit.l.failed,
            this.bit.l.updateError
          );
        }
      } else {
        this.notification.success(
          this.bit.l.success,
          this.bit.l.updateSuccess
        );
        this.getInformation();
      }
    });
  }
}
