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

  players: Player[];
  loading: boolean;

  constructor(private playersService: PlayersService) { }

  ngOnInit() {
    this.search.valueChanges
      .subscribe(value => {
        this.playersService.searchText = value;
        if (value === '') {
          this.getPlayers();
        } else {
          this.searchPlayers(value);
        }
      });
    this.search.patchValue(this.playersService.searchText);
  }

  searchPlayers(search: string): void {
    this.loading = true;
    this.playersService.searchPlayers(search)
      .subscribe( players => {
        this.players = players;
        this.loading = false;
      });
  }

  getPlayers(): void {
    this.loading = true;
    this.playersService.getPlayers()
      .subscribe( players => {
        this.players = players.values;
        this.loading = false;
      });
  }

  public clearSearch(): void {
    this.search.patchValue('');
  }

}
