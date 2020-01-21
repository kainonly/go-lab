import {en_US, zh_CN} from 'ng-zorro-antd';

export const environment = {
  production: true,
  bit: {
    originUrl: 'https://api.kainonly.com',
    staticUrl: 'https://cdn.kainonly.com/',
    iconUrl: 'https://cdn.kainonly.com/',
    namespace: '/system',
    uploadsUrl: false,
    uploadsPath: 'system/main/uploads',
    withCredentials: true,
    httpInterceptor: true,
    pageLimit: 10,
    breadcrumbTop: 0,
    formControlCol: {
      common: {
        nzXXl: 8,
        nzXl: 9,
        nzLg: 10,
        nzMd: 14,
        nzSm: 24,
      },
      submit: {
        nzXXl: {span: 8, offset: 4},
        nzXl: {span: 9, offset: 5},
        nzLg: {span: 10, offset: 6},
        nzMd: {span: 14, offset: 6},
        nzSm: {span: 24, offset: 0}
      }
    },
    formLabelCol: {
      common: {
        nzXXl: 4,
        nzXl: 5,
        nzLg: 6,
        nzMd: 7,
        nzSm: 24
      },
    },
    localDefault: 'zh_cn',
    localeBind: new Map([
      ['zh_cn', zh_CN],
      ['en_us', en_US]
    ]),
    i18nDefault: 'zh_cn',
    i18nContain: ['zh_cn', 'en_us'],
    i18nSwitch: [
      {
        i18n: 'zh_cn',
        name: {
          zh_cn: '中文',
          en_us: 'Chinese'
        }
      },
      {
        i18n: 'en_us',
        name: {
          zh_cn: '英文',
          en_us: 'English'
        }
      }
    ]
  }
};
