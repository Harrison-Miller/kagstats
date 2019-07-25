import { Component, OnDestroy, OnInit } from '@angular/core';
import { PlayersService } from './services/players.service';
import { takeUntil } from 'rxjs/operators';
import { Subject } from 'rxjs';
import { Player } from './models';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit, OnDestroy {
  componentDestroyed$ = new Subject();
  players$ = new Subject<Player[]>();

  constructor(public playersService: PlayersService) {}

  ngOnInit() {
    this.playersService
      .getPlayers()
      .pipe(takeUntil(this.componentDestroyed$))
      .subscribe(page => {
        this.players$.next(page.values);
      });
  }

  ngOnDestroy() {
    this.componentDestroyed$.next();
  }
}
