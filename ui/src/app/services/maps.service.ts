import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { MapBasics, GithubTree } from '../models';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class MapsService {

  constructor(private http: HttpClient) { }

  getMaps(): Observable<MapBasics[]> {
    let path = `${environment.apiUrl}/maps`;
    return this.http.get<MapBasics[]>(path);
  }

  getMapPaths(): Observable<GithubTree> {
    let path = "https://api.github.com/repos/transhumandesign/kag-base/git/trees/master?recursive=true";
    return this.http.get<GithubTree>(path);
  }
}
