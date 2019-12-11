import {CUSTOM_ELEMENTS_SCHEMA, NgModule} from '@angular/core';
import {NgZorroAntdModule} from 'ng-zorro-antd';
import {CommonModule} from '@angular/common';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {NgxBitDirectiveModule, NgxBitExtModule, NgxBitPipeModule} from 'ngx-bit';

@NgModule({
  exports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    NgZorroAntdModule,
    NgxBitExtModule,
    NgxBitPipeModule,
    NgxBitDirectiveModule,
  ],
  schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class AppExtModule {
}
