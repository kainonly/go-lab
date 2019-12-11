import {Component, OnInit} from '@angular/core';
import {AbstractControl, FormBuilder, FormGroup, Validators} from '@angular/forms';
import {SwalService, BitService, asyncValidator} from 'ngx-bit';
import {NzNotificationService} from 'ng-zorro-antd';
import {AclService} from '@common/acl.service';
import packer from './language';
import {switchMap} from 'rxjs/operators';

@Component({
  selector: 'app-acl-add',
  templateUrl: './acl-add.component.html'
})
export class AclAddComponent implements OnInit {
  form: FormGroup;
  defaultWriteLists: any[] = [
    {label: 'add', value: 'add'},
    {label: 'edit', value: 'edit'},
    {label: 'delete', value: 'delete'},
  ];
  writeLists: any[] = [];
  writeVisible = false;
  writeValue = '';
  defaultReadLists: any[] = [
    {label: 'get', value: 'get'},
    {label: 'originLists', value: 'originLists'},
    {label: 'lists', value: 'lists'},
  ];
  readLists: any[] = [];
  readVisible = false;
  readValue = '';

  constructor(
    public bit: BitService,
    private fb: FormBuilder,
    private notification: NzNotificationService,
    private swal: SwalService,
    private aclService: AclService
  ) {
  }

  ngOnInit() {
    this.bit.registerLocales(packer);
    this.form = this.fb.group({
      name: this.fb.group(this.bit.i18nGroup({
        validate: {
          zh_cn: [Validators.required],
          en_us: []
        }
      })),
      key: [null, [Validators.required], [this.existsKey]],
      status: [true, [Validators.required]]
    });
  }

  existsKey = (control: AbstractControl) => {
    return asyncValidator(this.aclService.validedKey(control.value));
  };

  /**
   * 显示新增写入表单
   */
  showAddWriteAction() {
    this.writeVisible = true;
  }

  /**
   * 提交新增写入表单
   */
  submitAddWriteAction() {
    if (this.writeLists.indexOf(this.writeValue) !== -1) {
      this.notification.error(
        this.bit.l.addFailed,
        this.bit.l.addOperateFailed
      );
      return;
    }

    this.writeLists.push(this.writeValue);
    this.writeValue = '';
    this.writeVisible = false;
  }

  /**
   * 取消新增写入表单
   */
  cancelAddWriteAction() {
    this.writeValue = '';
    this.writeVisible = false;
  }

  /**
   * 删除写入操作
   */
  deleteAddAction(operate: string) {
    this.writeLists.splice(
      this.writeLists.indexOf(operate),
      1
    );
  }

  /**
   * 显示新增读取表单
   */
  showAddReadAction() {
    this.readVisible = true;
  }

  /**
   * 提交新增读取表单
   */
  submitAddReadAction() {
    if (this.readLists.indexOf(this.readValue) !== -1) {
      this.notification.error(
        this.bit.l.addFailed,
        this.bit.l.addOperateFailed
      );
      return;
    }

    this.readLists.push(this.readValue);
    this.readValue = '';
    this.readVisible = false;
  }

  /**
   * 取消新增读取表单
   */
  cancelAddReadAction() {
    this.readValue = '';
    this.readVisible = false;
  }

  /**
   * 删除读取操作
   */
  deleteReadAction(operate: string) {
    this.readLists.splice(
      this.readLists.indexOf(operate),
      1
    );
  }

  /**
   * 提交
   */
  submit(data) {
    Reflect.set(data, 'name', JSON.stringify(data.name));
    Reflect.set(data, 'write', [
      ...this.defaultWriteLists.filter(v => v.checked).map(v => v.value),
      ...this.writeLists
    ].join(','));
    Reflect.set(data, 'read', [
      ...this.defaultReadLists.filter(v => v.checked).map(v => v.value),
      ...this.readLists
    ].join(','));
    this.aclService.add(data).pipe(
      switchMap(res => this.swal.addAlert(res, this.form, {
        status: true
      }))
    ).subscribe(() => {
    });
  }
}
