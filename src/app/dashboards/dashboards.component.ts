import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { BitService, BitEventsService, BitSupportService } from 'ngx-bit';
import { NzNotificationService } from 'ng-zorro-antd';
import { MainService } from '@common/main.service';

@Component({
  selector: 'app-dashboards',
  templateUrl: './dashboards.component.html',
  styleUrls: ['./dashboards.component.scss']
})
export class DashboardsComponent implements OnInit, OnDestroy {
  collapsed = false;
  navLists: any[] = [];

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private mainService: MainService,
    private events: BitEventsService,
    private notification: NzNotificationService,
    public support: BitSupportService,
    public bit: BitService
  ) {
  }

  ngOnInit() {
    this.getMenuLists();
    this.support.setup(this.router);
    this.events.on('locale').subscribe(locale => {
      this.bit.locale = locale;
    });
    this.events.on('refresh-menu').subscribe(() => {
      this.getMenuLists();
    });
  }

  ngOnDestroy() {
    this.events.off('locale');
    this.events.off('refresh-menu');
    this.support.destory();
  }

  /**
   * Get Menu Lists
   */
  private getMenuLists() {
    this.mainService.resource().subscribe(data => {
      this.support.putResource(data.resource, data.router);
      this.navLists = data.nav;
    });
  }

  /**
   * User logout
   */
  logout() {
    this.mainService.logout().subscribe(() => {
      this.support.clear();
      this.support.destory();
      this.router.navigateByUrl('/login');
      this.notification.success(this.bit.l.logout, this.bit.l.logoutSuccess);
    });
  }
}
