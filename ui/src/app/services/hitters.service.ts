import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Hitter } from '../models';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class HittersService {

  constructor(private http: HttpClient) { }

  getHitters(playerId: number): Observable<Hitter[]> {
    let path = `${environment.apiUrl}/players/${playerId}/hitters`;
    return this.http.get<{hitters:Hitter[]}>(path).pipe(
      map(h => h.hitters)
    );
  }
}
