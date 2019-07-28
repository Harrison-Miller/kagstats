import { Component, OnInit, Input } from '@angular/core';
import { KillsService } from 'src/app/services/kills.service';
import { Kill } from '../../models';
import { Observable, Subject, timer } from 'rxjs';
import { map, take } from 'rxjs/operators';

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
  styleUrls: ['./killfeed.component.sass']
})
export class KillfeedComponent implements OnInit {
  @Input() limit: number = 100;
  kills$: Observable<Kill[]>;

  constructor(private killsService: KillsService) {}

  ngOnInit() {
    this.kills$ = this.killsService.kills$.pipe(
      map(kills => {
        return kills.slice(0, this.limit);
      })
    );
  }
}
