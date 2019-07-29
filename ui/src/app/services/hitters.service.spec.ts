import { TestBed } from '@angular/core/testing';

import { HittersService } from './hitters.service';

describe('HittersService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: HittersService = TestBed.get(HittersService);
    expect(service).toBeTruthy();
  });
});
