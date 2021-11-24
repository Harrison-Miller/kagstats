import { Injectable } from '@angular/core';
import {ActivatedRouteSnapshot, CanActivate, RouterStateSnapshot, UrlTree} from '@angular/router';
import {Observable, of} from 'rxjs';
import {AuthService} from "../services/auth.service";
import {catchError, map} from "rxjs/operators";

@Injectable({
  providedIn: 'root'
})
export class WithPermissionGuard implements  CanActivate {
  constructor(private authService: AuthService) {
  }

  canActivate(next: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> {
    const permissionsRequired = next.data['permissionsRequired'];
    return this.authService.validate().pipe(map( resp => {
      const claims = this.authService.playerClaims.getValue();
      for (const p of permissionsRequired) {
        if (!claims.permissions.includes(p)) {
          return false;
        }
      }
      return true;
    }), catchError((err) => {
      return of(false);
    }));
  }
}
