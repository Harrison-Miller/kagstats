import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {PollAnswer, PollResponse} from "../models";
import {Observable} from 'rxjs';
import {environment} from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class PollService {

  constructor(private http: HttpClient) { }

  getCurrentPoll(): Observable<PollResponse> {
    const path = `${environment.apiUrl}/poll`;
    return this.http.get<PollResponse>(path);
  }

  answerCurrentPoll(answers: PollAnswer[]): Observable<any> {
    const path = `${environment.apiUrl}/poll`;
    return this.http.post<any>(path, answers);
  }
}

