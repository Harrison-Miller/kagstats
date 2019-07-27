import { Component, OnInit } from '@angular/core';
import { FormControl } from '@angular/forms';
import { PlayersService } from '../../services/players.service';
import { Player } from '../../models';

@Component({
  selector: 'app-players',
  templateUrl: './players.component.html',
  styleUrls: ['./players.component.scss']
})
export class PlayersComponent implements OnInit {

  search = new FormControl('')

  players: Player[]
  loading: boolean

  constructor(private playersService: PlayersService) { 
    this.getPlayers();
  }

  ngOnInit() {
    this.search.valueChanges
      .subscribe(value => {
        if(value == '') {
          this.getPlayers();
        }
        this.searchPlayers(value);
      })
  }

  searchPlayers(search: string): void {
    this.loading = true;
    this.playersService.searchPlayers(search)
      .subscribe( players => {
        this.players = players;

        if(this.players) {
          if(this.players.length != 0) {
            this.getAvatars();
            this.getAPIPlayers();
          }
        }
        this.loading = false;
      });
  }

  getPlayers(): void {
    this.loading = true;
    this.playersService.getPlayers()
      .subscribe( players => {
        this.players = players.values;
        this.getAvatars();
        this.getAPIPlayers();
        this.loading = false;
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
