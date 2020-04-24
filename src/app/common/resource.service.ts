import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import { HttpService } from 'ngx-bit';

@Injectable()
export class ResourceService {
  private model = 'resource';

  constructor(
    private http: HttpService
  ) {
  }

  originLists(): Observable<any> {
    return this.http.originLists(this.model);
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

  /**
   * Resource Sort
   */
  sort(data: any[]): Observable<any> {
    return this.http.req(this.model + '/sort', {
      data
    });
  }

  /**
   * Validate Resource Key
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


