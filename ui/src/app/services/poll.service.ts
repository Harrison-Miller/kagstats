import { Injectable } from '@angular/core';
import {HttpClient, HttpResponse} from "@angular/common/http";
import {PollAnswer, PollCompletedResponse, Poll} from "../models";
import {Observable} from 'rxjs';
import {environment} from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class PollService {

  constructor(private http: HttpClient) { }

  listPolls(): Observable<Poll[]> {
    const path = `${environment.apiUrl}/polls`;
    return this.http.get<Poll[]>(path);
  }

  pollCompleted(): Observable<PollCompletedResponse> {
    const path = `${environment.apiUrl}/poll/completed`;
    return this.http.get<PollCompletedResponse>(path);
  }

  getCurrentPoll(): Observable<Poll> {
    const path = `${environment.apiUrl}/poll`;
    return this.http.get<Poll>(path);
  }

  answerCurrentPoll(answers: PollAnswer[]): Observable<any> {
    const path = `${environment.apiUrl}/poll`;
    return this.http.post<any>(path, answers);
  }

  downloadPoll(pollID: number): Observable<HttpResponse<any>> {
    const path = `${environment.apiUrl}/poll/${pollID}/download`;
    return this.http.get<any>(path, { observe: 'response', responseType: 'blob' as 'json' });
  }
}

