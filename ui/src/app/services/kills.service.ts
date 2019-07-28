import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { PagedResult, Kill } from '../models';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class KillsService {

  constructor(private http: HttpClient) { }

  /**
   * Perform GET request at our kills endpoint. 
   * We can optionally provide a start parameter in our request, 
   * change the starting point of 
   * 
   * @param start Optional starting point for paged results.
   * 
   * @todo Ask for clarification about whether optional params are okay or not.
   */
  getKills(start?: number): Observable<PagedResult<Kill>> {

    // Only create a parameter object if passed a start value. 
    const options = start ? 
      { params: new HttpParams().set('start', start.toString()) } : {};
    
    return this.http.get<PagedResult<Kill>>('/api/kills', options);
  }

}
