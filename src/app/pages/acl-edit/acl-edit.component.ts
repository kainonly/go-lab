import {Component, OnInit} from '@angular/core';
import {SwalService, BitService, asyncValidator} from 'ngx-bit';
import {AbstractControl, AsyncValidatorFn, FormBuilder, FormGroup, Validators} from '@angular/forms';
import {NzNotificationService} from 'ng-zorro-antd';
import {map, switchMap} from 'rxjs/operators';
import {AclService} from '@common/acl.service';
import packer from './language';
import {ActivatedRoute} from '@angular/router';
import {AsyncSubject} from 'rxjs';

@Component({
  selector: 'app-acl-edit',
  templateUrl: './acl-edit.component.html'
})
export class AclEditComponent implements OnInit {
  private id: number;
  private nameAsync: AsyncSubject<string> = new AsyncSubject<string>();
  private keyAsync: AsyncSubject<string> = new AsyncSubject();
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
    private route: ActivatedRoute,
    private aclService: AclService,
  ) {
  }

  ngOnInit() {
    this.bit.registerLocales(packer);
    this.form = this.fb.group({
      name: this.fb.group(this.bit.i18nGroup({
        validate: {
          zh_cn: [Validators.required],
          en_us: []
        },
        asyncValidate: {
          zh_cn: [this.existsName]
        }
      })),
      key: [null, [Validators.required], [this.existsKey]],
      status: [true, [Validators.required]]
    });
    this.route.params.subscribe(param => {
      this.id = param.id;
      this.getData();
    });
  }

  existsName: AsyncValidatorFn = (control: AbstractControl) => {
    return asyncValidator(this.aclService.validedName(control.value, this.nameAsync)).pipe(
      map(result => {
        if (control.touched) {
          control.setErrors(result);
        }
        return result;
      })
    );
  };

  existsKey = (control: AbstractControl) => {
    return asyncValidator(this.aclService.validedKey(control.value, this.keyAsync));
  };

  getData() {
    this.aclService.get(this.id).subscribe(data => {
      if (!data) {
        return;
      }
      const name = this.bit.i18nParse(data.name);
      this.nameAsync.next(name.zh_cn);
      this.nameAsync.complete();
      this.keyAsync.next(data.key);
      this.keyAsync.complete();
      if (data.write) {
        data.write.split(',').forEach((value) => {
          if (['add', 'edit', 'delete'].indexOf(value) !== -1) {
            this.defaultWriteLists.filter(v => value === v.value).map(v => {
              v.checked = true;
              return v;
            });
          } else {
            this.writeLists.push(value);
          }
        });
      }
      if (data.read) {
        data.read.split(',').forEach((value) => {
          if (['get', 'originLists', 'lists'].indexOf(value) !== -1) {
            this.defaultReadLists.filter(v => value === v.value).map(v => {
              v.checked = true;
              return v;
            });
          } else {
            this.readLists.push(value);
          }
        });
      }
      this.form.setValue({
        name,
        key: data.key,
        status: data.status
      });
    });
  }

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
    Reflect.set(data, 'id', this.id);
    Reflect.set(data, 'name', JSON.stringify(data.name));
    Reflect.set(data, 'write', [
      ...this.defaultWriteLists.filter(v => v.checked).map(v => v.value),
      ...this.writeLists
    ].join(','));
    Reflect.set(data, 'read', [
      ...this.defaultReadLists.filter(v => v.checked).map(v => v.value),
      ...this.readLists
    ].join(','));
    this.aclService.edit(data).pipe(
      switchMap(res => this.swal.editAlert(res))
    ).subscribe((status) => {
      if (status) {
        this.getData();
      }
    });
  }
}
