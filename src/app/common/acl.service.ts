import {Injectable} from '@angular/core';
import {Observable, of} from 'rxjs';
import {switchMap} from 'rxjs/operators';
import {HttpService} from 'ngx-bit';

@Injectable()
export class AclService {
  private model = 'acl';

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

  add(data: any) {
    return this.http.add(this.model, data);
  }

  get(id: number) {
    return this.http.get(this.model, id);
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
   * Validate Acl Key
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
