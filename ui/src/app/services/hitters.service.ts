import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import {Hitter, MonthlyHittersStats} from '../models';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class HittersService {

  constructor(private http: HttpClient) { }

  getHitters(playerID: number): Observable<Hitter[]> {
    const path = `${environment.apiUrl}/players/${playerID}/hitters`;
    return this.http.get<{hitters: Hitter[]}>(path).pipe(
      map(h => h.hitters)
    );
  }

  getMonthlyHitters(playerID: number): Observable<MonthlyHittersStats[]> {
    const path = `${environment.apiUrl}/players/${playerID}/hitters/monthly`;
    return this.http.get<MonthlyHittersStats[]>(path);
  }
}
