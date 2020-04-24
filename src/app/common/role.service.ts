import { Injectable } from '@angular/core';
import { HttpService } from 'ngx-bit';
import { Observable, of } from 'rxjs';
import { switchMap } from 'rxjs/operators';

@Injectable()
export class RoleService {
  private model = 'role';

  constructor(
    private http: HttpService
  ) {
  }

  originLists(): Observable<any> {
    return this.http.originLists(this.model);
  }

  lists(search: any, refresh: boolean): Observable<any> {
    return this.http.lists(this.model, search, refresh);
  }

  get(id: number) {
    return this.http.get(this.model, id);
  }

  add(data: any) {
    return this.http.add(this.model, data);
  }

  edit(data: any): Observable<any> {
    return this.http.edit(this.model, data);
  }

  delete(id: number | number[]): Observable<any> {
    return this.http.delete(this.model, Array.isArray(id) ? id : [id]);
  }

  status(data: any): Observable<any> {
    return this.http.status(this.model, data);
  }

  /**
   * Validate Role Key
   */
  validedKey(key: string, edit: Observable<string> = of(null)) {
    return edit.pipe(
      switchMap(editKey => (key === editKey ? of({
          error: 0,
          data: false
        }) : this.http.req(this.model + '/validedKey', {
          key
        })
      ))
    );
  }
}
