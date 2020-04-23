import {Injectable} from '@angular/core';
import {HttpService} from 'ngx-bit';
import {Observable} from 'rxjs';

@Injectable()
export class PolicyService {
  private model = 'policy';

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

  delete(id: number | number[]): Observable<any> {
    return this.http.delete(this.model, Array.isArray(id) ? id : [id]);
  }
}