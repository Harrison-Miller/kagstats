import { Component, OnInit, Input } from '@angular/core';
import { Player } from '../../models';

@Component({
  selector: 'app-player-banner',
  templateUrl: './player-banner.component.html',
  styleUrls: ['./player-banner.component.scss']
})
export class PlayerBannerComponent implements OnInit {

  @Input() player: Player;
  today: number = Date.now();

  // optional inputs for changing the style of the player banner
  @Input() avatar: boolean = true; // show the avatar
  @Input() status: boolean =  true; //show the status icon
  @Input() accolades: boolean = true; // show accolade list

  constructor() { 

  }

  ngOnInit() {
  }

}
