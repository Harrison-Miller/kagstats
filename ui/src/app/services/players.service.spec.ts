import { TestBed } from '@angular/core/testing';

import { PlayersService } from './players.service';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('PlayersService', () => {
  beforeEach(() => TestBed.configureTestingModule({
    imports: [HttpClientTestingModule],
    providers: [PlayersService]
  }));

  it('should be created', () => {
    const service: PlayersService = TestBed.get(PlayersService);
    expect(service).toBeTruthy();
  });
});
