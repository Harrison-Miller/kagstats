import { Component, OnDestroy, OnInit } from '@angular/core';
import { LeaderboardService } from '../../../services/leaderboard.service';
import { Subject } from 'rxjs';
import { BasicStats } from '../../../models';
import { takeUntil } from 'rxjs/operators';
import { ActivatedRoute, Router } from '@angular/router';

const availableClasses = ['Archer', 'Builder', 'Knight', 'Kills', 'MonthlyArcher', 'MonthlyBuilder', 'MonthlyKnight'];
const boardTitle = {
  'Archer': 'All Time Archer',
  'Builder': 'All Time Builder',
  'Knight': 'All Time Knight',
  'MonthlyArcher': 'Monthly Archer',
  'MonthlyBuilder': 'Monthly Builder',
  'MonthlyKnight': 'Monthly Knight',
  'Kills': 'Kills'
}

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
          this.router.navigateByUrl('/MonthlyArcher');
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
    if(this.class == "Kills") {
      return leader.totalKills;
    }
    return leader[`${this.class.toLowerCase().replace("monthly", "")}Kills`]
  }

  boardTitle(): string {
    return boardTitle[this.class];
  }

  boardDate(): string {
    if(this.class.toLocaleLowerCase().includes("monthly")) {

    }
    return "";
  }

  totalDeaths(leader: BasicStats): number {
    if(this.class == "Kills") {
      return leader.totalDeaths;
    }
    return leader[`${this.class.toLowerCase().replace("monthly", "")}Deaths`]
  }
  
  kd(leader: BasicStats): string {
    const deaths = this.totalDeaths(leader);
    return (this.totalKills(leader) / (deaths === 0 ? 1 : deaths)).toFixed(2);
  }
}
