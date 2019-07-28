import { Component, OnDestroy, OnInit } from '@angular/core';
import { LeaderboardService } from '../../../services/leaderboard.service';
import { Subject } from 'rxjs';
import { BasicStats } from '../../../models';
import { takeUntil } from 'rxjs/operators';
import { ActivatedRoute, Router } from '@angular/router';

const availableClasses = ['Archer', 'Builder', 'Knight'];

@Component({
  selector: 'app-class-leaderboard',
  templateUrl: './class-leaderboard.component.html',
  styleUrls: ['./class-leaderboard.component.sass']
})
export class ClassLeaderboardComponent implements OnInit, OnDestroy {
  class: string;
  componentDestroyed$ = new Subject();
  leaderboard$ = new Subject<BasicStats[]>();

  constructor(
    private leaderboardService: LeaderboardService,
    private route: ActivatedRoute,
    private router: Router
  ) {}

  ngOnInit() {
    this.route.paramMap
      .pipe(takeUntil(this.componentDestroyed$))
      .subscribe(params => {
        this.class = params.get('board');

        if (!availableClasses.includes(this.class)) {
          this.router.navigateByUrl('/leaderboards');
          console.error('Invalid class');
        }

        this.leaderboardService
          .leaderboard
          .pipe(takeUntil(this.componentDestroyed$))
          .subscribe(result => {
            this.leaderboard$.next(result.leaderboard);
          });

        this.leaderboardService.getLeaderboard(this.class);
      });
  }

  ngOnDestroy() {
    this.componentDestroyed$.next();
  }

  totalKills(leader: BasicStats): number {
    return leader[`${this.class.toLowerCase()}Kills`]
  }

  totalDeaths(leader: BasicStats): number {
    return leader[`${this.class.toLowerCase()}Deaths`]
  }

  kd(leader: BasicStats): string {
    const c = this.class.toLowerCase();
    return (
      leader[`${c}Kills`] / (leader[`${c}Deaths`] === 0 ? 1 : leader[`${c}Deaths`])
    ).toFixed(2);
  }
}
