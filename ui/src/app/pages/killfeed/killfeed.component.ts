import { Component, OnInit } from '@angular/core';
import { KillsService } from 'src/app/services/kills.service';
import { Kill } from '../../models';
import { Observable, timer } from 'rxjs';

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
  kills: Kill[];

  constructor(private killsService: KillsService) {}

  ngOnInit() {
    // Routinely query for new kills. 
    timer(0, 30000).subscribe(() => {
      this.killsService.getKills().subscribe(kills => {
        // Map kills to values prop of returned PagedResults.
        this.kills = kills.values;
      });
    });
  }
}
