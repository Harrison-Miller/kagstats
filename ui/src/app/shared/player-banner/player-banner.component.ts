import { Component, OnInit, Input } from '@angular/core';
import { Player } from '../../models';

export const REGISTERED_TOOLTIP = {
  "r10yBadge": "Playing for over 10 years",
  "r9yBadge": "Playing for over 9 years",
  "r8yBadge": "Playing for over 8 years",
  "r7yBadge": "Playing for over 7 years",
  "r6yBadge": "Playing for over 6 years",
  "r5yBadge": "Playing for over 5 years",
  "r4yBadge": "Playing for over 4 years",
  "r3yBadge": "Playing for over 3 years",
  "r2yBadge": "Playing for over 2 years",
  "r1yBadge": "Playing for 1 year",
  "r9mBadge": "Playing for 9 months",
  "r6mBadge": "Playing for 6 months",
  "r3mBadge": "Playing for 3 months",
  "r2mBadge": "Playing for 2 months",
  "r1mBadge": "Playing for 1 month",
  "r3wBadge": "Playing for 3 weeks",
  "r2wBadge": "Playing for 2 weeks",
  "r1wBadge": "Playing for 1 week",
  "newBadge": "Welcome this player to the game!"
};

@Component({
  selector: 'app-player-banner',
  templateUrl: './player-banner.component.html',
  styleUrls: ['./player-banner.component.scss']
})
export class PlayerBannerComponent implements OnInit {

  public registered_tooltip = REGISTERED_TOOLTIP

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
