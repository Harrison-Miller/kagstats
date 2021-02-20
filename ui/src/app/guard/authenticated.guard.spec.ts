import { TestBed, async, inject } from '@angular/core/testing';

import { AuthenticatedGuard } from './authenticated.guard';

describe('AuthenticatedGuard', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [AuthenticatedGuard]
    });
  });

  it('should ...', inject([AuthenticatedGuard], (guard: AuthenticatedGuard) => {
    expect(guard).toBeTruthy();
  }));
});
