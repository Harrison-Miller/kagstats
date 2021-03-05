import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { PagedResult, Kill } from '../models';
import {BehaviorSubject, Observable, Subject, timer} from 'rxjs';
import { environment } from '../../environments/environment';
import {takeUntil} from "rxjs/operators";

@Injectable({
  providedIn: 'root'
})
export class KillsService {
  kills = new BehaviorSubject<Kill[]>(null);

  constructor(private http: HttpClient) {
  }

  nextGetKills$ = new Subject();

  getKills(url: string, start?: number, limit?: number) {
    this.nextGetKills$.next(null);

    // Only create a parameter object if passed a start value.
    const options = { params: new HttpParams() };
    if (start) {
      options.params = options.params.append('start', start.toString());
    }

    if (limit) {
      options.params = options.params.append('limit', limit.toString());
    }

    this.http.get<PagedResult<Kill>>(environment.apiUrl + url, options).subscribe(result => {
      this.kills.next(result.values);
    });

    timer(0, 10000).pipe(takeUntil(this.nextGetKills$))
      .subscribe(x => {
        this.http.get<PagedResult<Kill>>(environment.apiUrl + url, options).subscribe(result => {
          this.kills.next(result.values);
        });
    });
  }
}
