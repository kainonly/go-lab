import { Component, OnInit } from '@angular/core';
import { AbstractControl, AsyncValidatorFn, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { SwalService, BitService, asyncValidator } from 'ngx-bit';
import { NzNotificationService } from 'ng-zorro-antd';
import { AclService } from '@common/acl.service';
import { map, switchMap } from 'rxjs/operators';

@Component({
  selector: 'app-acl-add',
  templateUrl: './acl-add.component.html'
})
export class AclAddComponent implements OnInit {
  form: FormGroup;
  writeLists: string[] = ['add', 'edit', 'delete'];
  readLists: string[] = ['get', 'originLists', 'lists'];

  constructor(
    public bit: BitService,
    private fb: FormBuilder,
    private notification: NzNotificationService,
    private swal: SwalService,
    private aclService: AclService
  ) {
  }

  ngOnInit() {
    this.bit.registerLocales(import('./language'));
    this.form = this.fb.group({
      name: this.fb.group(this.bit.i18nGroup({
        validate: {
          zh_cn: [Validators.required],
          en_us: []
        },
        asyncValidate: {
          zh_cn: [this.existsName],
          en_us: []
        }
      })),
      key: [null, [Validators.required], [this.existsKey]],
      write: [this.writeLists],
      read: [this.readLists],
      status: [true, [Validators.required]]
    });
  }

  existsName: AsyncValidatorFn = (control: AbstractControl) => {
    return asyncValidator(this.aclService.validedName(control.value)).pipe(
      map(result => {
        if (control.touched) {
          control.setErrors(result);
        }
        return result;
      })
    );
  };

  existsKey = (control: AbstractControl) => {
    return asyncValidator(this.aclService.validedKey(control.value));
  };

  /**
   * 提交
   */
  submit(data) {
    this.aclService.add(data).pipe(
      switchMap(res => this.swal.addAlert(res, this.form, {
        status: true
      }))
    ).subscribe(() => {
    });
  }
}
