import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { PagedResult, Kill } from '../models';
import { BehaviorSubject, Observable, timer } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class KillsService {
  kills$ = new BehaviorSubject<Kill[]>([]);

  constructor(private http: HttpClient) {
    timer(0, 5000).subscribe(() => {
      this.getKills();
    });
  }

  /**
   * Perform GET request at our kills endpoint.
   * We can optionally provide a start parameter in our request,
   * change the starting point of
   *
   * @param start Optional starting point for paged results.
   *
   * @todo Ask for clarification about whether optional params are okay or not.
   */
  getKills(start?: number) {
    // Only create a parameter object if passed a start value.
    const options = start
      ? { params: new HttpParams().set('start', start.toString()) }
      : {};

    this.http.get<PagedResult<Kill>>('/api/kills', options).subscribe(result => {
      this.kills$.next(result.values);
    })
  }
}
