import { Component, OnInit } from '@angular/core';
import { PlayersService } from '../../services/players.service';
import { Player } from '../../models';

@Component({
  selector: 'app-players',
  templateUrl: './players.component.html',
  styleUrls: ['./players.component.scss']
})
export class PlayersComponent implements OnInit {

  players: Player[]

  constructor(private playersService: PlayersService) { 
    this.getPlayers();
  }

  ngOnInit() {
  }

  getPlayers(): void {
    this.playersService.getPlayers()
      .subscribe( players => {
        this.players = players.values
        this.getAvatars();
        this.getAPIPlayers();
      });
  }

  getAvatars(): void {
    this.players.forEach(player => {
      this.playersService.getPlayerAvatar(player.username)
        .subscribe(avatar => player.avatar = avatar)
    });
  }
  
  getAPIPlayers(): void {
    this.players.forEach(player => {
      this.playersService.getAPIPlayer(player.username)
        .subscribe(p => {
          player.info = p.playerInfo;
          player.status = p.playerStatus;
        })
    });
  }

}
