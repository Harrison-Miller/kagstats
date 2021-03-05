import {Component, OnInit, Input, OnChanges, OnDestroy} from '@angular/core';
import { KillsService } from 'src/app/services/kills.service';
import { Kill } from '../../models';
import { Observable, Subject, timer } from 'rxjs';
import {map, take, takeUntil} from 'rxjs/operators';
import { HITTER_DESCRIPTION } from '../../hitters';

/**
 * The KillfeedComponent periodically queries for new kills,
 * and displays a list detailing the recent most kills.
 *
 *
 * @todo Consider modularizing component by having it accept an array of kills, and performing the query elsewhere.
 */
@Component({
  selector: 'app-killfeed',
  templateUrl: './killfeed.component.html',
  styleUrls: ['./killfeed.component.scss']
})
export class KillfeedComponent implements OnInit, OnDestroy {
  @Input() limit = 100;
  @Input() url = '/kills';
  kills: Kill[];
  descriptions: string[] = HITTER_DESCRIPTION;

  componentDestroyed$ = new Subject();
  loading = true;

  constructor(private killsService: KillsService) {}

  ngOnInit() {
    this.killsService.kills.next(null);
    this.killsService.kills.pipe(takeUntil(this.componentDestroyed$))
      .subscribe( value => {
        if (value) {
          this.loading = false;
          this.kills = value;
        }
    });
    this.loading = true;
    this.killsService.getKills(this.url, 0, this.limit);
  }

  ngOnDestroy(): void {
    this.componentDestroyed$.next();
    this.killsService.nextGetKills$.next(null);
  }
}
