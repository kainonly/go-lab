import { Component, OnInit } from '@angular/core';
import { BitService } from 'ngx-bit';

@Component({
  selector: 'app-empty',
  template: `
    <nz-result nzStatus="403" nzTitle="403" [nzSubTitle]="bit.l['tips']">
      <div nz-result-extra>
        <button nz-button
                nzType="primary"
                routerLink="/">
          {{bit.l['home']}}
        </button>
      </div>
    </nz-result>
  `
})
export class EmptyComponent implements OnInit {
  constructor(
    public bit: BitService
  ) {
  }

  ngOnInit() {
    this.bit.registerLocales({
      home: ['回到主页', 'Go Dashboard'],
      tips: ['抱歉，您访问的页面不存在或未授权', 'Sorry, the page you visited does not exist or is not authorized']
    });
  }
}
