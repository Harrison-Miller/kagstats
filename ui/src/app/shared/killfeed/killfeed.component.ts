import { Component, OnInit, Input, OnChanges } from '@angular/core';
import { KillsService } from 'src/app/services/kills.service';
import { Kill } from '../../models';
import { Observable, Subject, timer } from 'rxjs';
import { map, take } from 'rxjs/operators';
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
export class KillfeedComponent implements OnInit, OnChanges {
  @Input() limit: number = 100;
  @Input() url: string = '/kills';
  kills$: Observable<Kill[]>;
  descriptions: string[] = HITTER_DESCRIPTION;

  constructor(private killsService: KillsService) {}

  ngOnInit() {
    this.kills$ = this.killsService.kills$.pipe(
      map(kills => {
        if(kills) {
          return kills.slice(0, this.limit);
        }
      })
    );
  }

  ngOnChanges() {
    this.killsService.url$.next(this.url);
  }
}
