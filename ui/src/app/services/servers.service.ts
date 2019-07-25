import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Server } from '../models';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ServersService {
  constructor(private http: HttpClient) {}

  getServers(): Observable<Server[]> {
    return this.http.get<Server[]>('/api/servers');
  }
}
