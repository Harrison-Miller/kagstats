import { TestBed } from '@angular/core/testing';

import { KillsService } from './kills.service';

describe('KillsService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: KillsService = TestBed.get(KillsService);
    expect(service).toBeTruthy();
  });
});
