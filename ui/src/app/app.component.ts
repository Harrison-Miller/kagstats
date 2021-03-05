import {Component, OnDestroy, OnInit} from '@angular/core';
import { Router } from '@angular/router';
import {NgbModal, ModalDismissReasons} from '@ng-bootstrap/ng-bootstrap';
import { FormBuilder } from '@angular/forms';
import {PlayersService} from './services/players.service';
import {AuthService} from './services/auth.service';
import {PlayerClaims} from './models';
import {Subject} from "rxjs";
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnDestroy, OnInit{

  componentDestroyed$ = new Subject();

  isNavbarCollapsed = true;

  closeResult = '';
  hideError = true;
  modal = null;

  playerClaims: PlayerClaims = null;

  constructor(public router: Router,
              private modalService: NgbModal,
              private formBuilder: FormBuilder,
              private authService: AuthService) {}

  loginForm = this.formBuilder.group({
    username: '',
    token: '',
  });

  ngOnInit(): void {
    this.authService.playerClaims.pipe(takeUntil(this.componentDestroyed$))
      .subscribe( value => {
        this.playerClaims = value;
      });
  }

  ngOnDestroy(): void {
    this.componentDestroyed$.next();
  }


  open(content) {
    this.modal = this.modalService.open(content, {ariaLabelledBy: 'modal-basic-title'});
    this.modal.result.then((result) => {
      this.loginForm.reset();
    }, (reason) => {
      this.loginForm.reset();
    });
  }

  private getDismissReason(reason: any): string {
    return 'closed';
  }

  generateToken(): void {
    const username = this.loginForm.value.username;
    window.open(`https://api.kag2d.com/v1/player/${username}/token/new`);
  }

  onSubmit(): void {
    const username = this.loginForm.value.username;
    const token = this.loginForm.value.token;
    this.authService.login(username, token).subscribe( noerr => {
      this.modal.close();
      this.loginForm.reset();
    },
    error => {
      this.hideError = false;
    });
  }

  logout(): void {
    this.authService.logout();
    this.router.navigate(['/leaderboards/MonthlyArcher']);
  }
}
