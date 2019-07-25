import { Component, OnDestroy, OnInit } from '@angular/core';
import { LeaderboardService } from '../../../services/leaderboard.service';
import { Subject } from 'rxjs';
import { BasicStats } from '../../../models';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-base-leaderboard',
  templateUrl: './base-leaderboard.component.html',
  styleUrls: ['./base-leaderboard.component.sass']
})
export class BaseLeaderboardComponent implements OnInit, OnDestroy {
  componentDestroyed$ = new Subject();
  leaderboard$ = new Subject<BasicStats[]>();

  constructor(private leaderboardService: LeaderboardService) {}

  ngOnInit() {
    this.leaderboardService
      .getBaseLeaderboard()
      .pipe(takeUntil(this.componentDestroyed$))
      .subscribe(result => {
        this.leaderboard$.next(result.leaderboard);
      });
  }

  ngOnDestroy() {
    this.componentDestroyed$.next();
  }

  totalKills(leader: BasicStats): number {
    return leader.archerKills + leader.builderKills + leader.knightKills;
  }

  totalDeaths(leader: BasicStats): number {
    return leader.archerDeaths + leader.builderDeaths + leader.knightDeaths;
  }

  kd(leader: BasicStats): Number {
    const deaths = this.totalDeaths(leader);
    return this.totalKills(leader) / (deaths === 0 ? 1 : deaths);
  }
}
