import { Component, OnInit } from '@angular/core';
import { BitConfigService, BitService } from 'ngx-bit';

@Component({
  selector: 'app-root',
  template: `
    <router-outlet></router-outlet>
  `
})
export class AppComponent implements OnInit {
  constructor(
    private bit: BitService,
    private config: BitConfigService
  ) {
  }

  ngOnInit() {
    this.config.setupLocales(import('./app.language'));
  }
}
