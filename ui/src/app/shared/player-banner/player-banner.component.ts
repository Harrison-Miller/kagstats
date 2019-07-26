import { Component, OnInit, Input } from '@angular/core';
import { Player } from '../../models';

@Component({
  selector: 'app-player-banner',
  templateUrl: './player-banner.component.html',
  styleUrls: ['./player-banner.component.scss']
})
export class PlayerBannerComponent implements OnInit {

  @Input() player: Player;

  constructor() { 

  }

  ngOnInit() {
  }
  

}
