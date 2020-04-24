import { Component, OnInit } from '@angular/core';
import { SwalService, BitService } from 'ngx-bit';
import { NzNotificationService } from 'ng-zorro-antd';
import { AclService } from '@common/acl.service';
import packer from './language';

@Component({
  selector: 'app-acl-index',
  templateUrl: './acl-index.component.html'
})
export class AclIndexComponent implements OnInit {
  lists = [];

  constructor(
    private swal: SwalService,
    public bit: BitService,
    public aclService: AclService,
    private notification: NzNotificationService
  ) {
  }

  ngOnInit() {
    this.bit.registerLocales(packer);
    this.bit.registerSearch('acl-index',
      { field: 'name->zh_cn', op: 'like', value: '' },
      { field: 'name->en_us', op: 'like', value: '' }
    ).subscribe(() => {
      this.getLists();
    });
  }

  /**
   * 获取列表数据
   */
  getLists(refresh = false) {
    this.aclService.lists(this.bit.getSearch(), refresh).subscribe(data => {
      this.lists = data;
    });
  }

  /**
   * 选择标签
   */
  selectType(value: number) {
    this.bit.search[0].value = value;
    this.getLists(true);
  }

  /**
   * 删除单操作
   */
  deleteData(id: any) {
    this.swal.deleteAlert(this.aclService.delete(id)).subscribe(res => {
      if (!res.error) {
        this.notification.success(
          this.bit.l.operateSuccess,
          this.bit.l.deleteSuccess
        );
        this.getLists(true);
      } else {
        this.notification.error(
          this.bit.l.operateError,
          this.bit.l.deleteError
        );
      }
    });
  }

  /**
   * 选中删除
   */
  deleteCheckData() {
    const id = this.lists.filter(value => value.checked).map(v => v.id);
    this.deleteData(id);
  }
}
