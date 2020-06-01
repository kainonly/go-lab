import { Component, OnInit } from '@angular/core';
import { BitService } from 'ngx-bit';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { NzNotificationService } from 'ng-zorro-antd';
import { MainService } from '@common/main.service';
import packer from './language';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html'
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
      old_password: [null, [this.validedPassword]],
      new_password: [null, [this.validedNewPassword]],
      new_password_check: [null, [this.checkNewPassword]]
    });
    this.getInformation();
  }

  validedPassword = (control: AbstractControl) => {
    if (!this.editPassword) {
      return;
    }
    if (control.parent === undefined) {
      return;
    }
    if (!control.value) {
      return { required: true };
    }
    const value = control.value;
    const len = value.length;
    if (len < 12) {
      return { min: true, error: true };
    }
    if (len > 20) {
      return { max: true, error: true };
    }
    if (value.match(/^(?=.*[a-z])[\w|@$!%*?&-+]+$/) === null) {
      return { lowercase: true, error: true };
    }
    if (value.match(/^(?=.*[A-Z])[\w|@$!%*?&-+]+$/) === null) {
      return { uppercase: true, error: true };
    }
    if (value.match(/^(?=.*[0-9])[\w|@$!%*?&-+]+$/) === null) {
      return { number: true, error: true };
    }
    if (value.match(/^(?=.*[@$!%*?&-+])[\w|@$!%*?&-+]+$/) === null) {
      return { symbol: true, error: true };
    }
    return null;
  };

  validedNewPassword = (control: AbstractControl) => {
    if (!this.editPassword) {
      return;
    }
    if (control.parent === undefined) {
      return;
    }
    if (!control.value) {
      return { required: true };
    }
    control.parent.get('new_password_check').updateValueAndValidity();
    const value = control.value;
    const len = value.length;
    if (len < 12) {
      return { min: true, error: true };
    }
    if (len > 20) {
      return { max: true, error: true };
    }
    if (value.match(/^(?=.*[a-z])[\w|@$!%*?&-+]+$/) === null) {
      return { lowercase: true, error: true };
    }
    if (value.match(/^(?=.*[A-Z])[\w|@$!%*?&-+]+$/) === null) {
      return { uppercase: true, error: true };
    }
    if (value.match(/^(?=.*[0-9])[\w|@$!%*?&-+]+$/) === null) {
      return { number: true, error: true };
    }
    if (value.match(/^(?=.*[@$!%*?&-+])[\w|@$!%*?&-+]+$/) === null) {
      return { symbol: true, error: true };
    }
    return null;
  };

  checkNewPassword = (control: AbstractControl) => {
    if (!this.editPassword) {
      return;
    }
    if (control.parent === undefined) {
      return;
    }
    if (!control.value) {
      return { required: true };
    }
    const password = control.parent.get('new_password').value;
    if (control.value !== password) {
      return { correctly: true, error: true };
    }
    return null;
  };

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
