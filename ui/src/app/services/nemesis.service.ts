import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Nemesis } from '../models';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class NemesisService {

  constructor(private http: HttpClient) { }

  getNemesis(playerId: number): Observable<Nemesis> {
    let path = `/api/players/${playerId}/nemesis`;
    return this.http.get<Nemesis>(path);
  }

  getBullied(playerId: number): Observable<Nemesis[]> {
    let path = `/api/players/${playerId}/bullied`;
    return this.http.get<{bullied:Nemesis[]}>(path).pipe(
      map(n => n.bullied)
    );
  }
}
