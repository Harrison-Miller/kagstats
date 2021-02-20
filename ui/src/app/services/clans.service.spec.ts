import { TestBed } from '@angular/core/testing';

import { ClansService } from './clans.service';

describe('ClansService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: ClansService = TestBed.get(ClansService);
    expect(service).toBeTruthy();
  });
});
