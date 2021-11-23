import {Component, OnDestroy, OnInit} from '@angular/core';
import {PollService} from "../../services/poll.service";
import {PlayerClaims, Poll} from "../../models";
import {AuthService} from "../../services/auth.service";
import {takeUntil} from "rxjs/operators";
import {Subject} from "rxjs";

@Component({
  selector: 'app-poll',
  templateUrl: './poll.component.html',
  styleUrls: ['./poll.component.sass']
})
export class PollComponent implements OnInit, OnDestroy {
  poll: Poll;
  playerClaims: PlayerClaims;
  componentDestroyed$ = new Subject();

  constructor(private pollService: PollService,
              private authService: AuthService) { }

  ngOnInit() {
    this.pollService.getCurrentPoll().subscribe(
      poll => {
        this.poll = poll;
        for(var q of this.poll.questions) {
          q.options_split = q.options.split(';');
        }
      }
    );

    this.authService.playerClaims.pipe(takeUntil(this.componentDestroyed$))
      .subscribe( value => {
        this.playerClaims = value;
      });
  }

  ngOnDestroy() {
    this.componentDestroyed$.next();
  }

}
