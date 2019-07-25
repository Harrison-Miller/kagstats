import { Component, OnDestroy, OnInit } from '@angular/core';
import { PlayersService } from './services/players.service';
import { takeUntil } from 'rxjs/operators';
import { Subject } from 'rxjs';
import { Player } from './models';
import { Router } from '@angular/router';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {

  constructor(public router: Router) {}
}
