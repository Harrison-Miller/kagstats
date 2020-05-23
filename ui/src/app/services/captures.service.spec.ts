import { TestBed } from '@angular/core/testing';

import { CapturesService } from './captures.service';

describe('CapturesService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: CapturesService = TestBed.get(CapturesService);
    expect(service).toBeTruthy();
  });
});
