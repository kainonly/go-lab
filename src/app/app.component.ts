import { Component, ComponentFactoryResolver, getPlatform, OnInit } from '@angular/core';
import { BitService, ConfigService } from 'ngx-bit';
import { NzMessageService } from 'ng-zorro-antd';
import { Observable, of } from 'rxjs';
import packer from './app.language';

@Component({
  selector: 'app-root',
  template: `
    <router-outlet></router-outlet>
  `
})
export class AppComponent implements OnInit {
  constructor(
    private bit: BitService,
    private message: NzMessageService,
    private config: ConfigService
  ) {
  }

  ngOnInit() {
    this.bit.registerLocales(packer, true);
    this.config.interceptor = (res: any): Observable<any> => {
      return of(res);
    };
  }
}
