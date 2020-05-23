import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Captures } from '../models';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class CapturesService {

  constructor(private http: HttpClient) { }

  getCaptures(playerId: number): Observable<Captures> {
    let path = `${environment.apiUrl}/players/${playerId}/captures`;
    return this.http.get<Captures>(path);
  }
}
