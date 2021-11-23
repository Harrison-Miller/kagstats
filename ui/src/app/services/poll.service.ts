import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Poll} from "../models";
import {Observable} from 'rxjs';
import {environment} from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class PollService {

  constructor(private http: HttpClient) { }

  getCurrentPoll(): Observable<Poll> {
    const path = `${environment.apiUrl}/poll`;
    return this.http.get<Poll>(path);
  }
}

