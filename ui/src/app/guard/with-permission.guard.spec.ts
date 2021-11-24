import { TestBed, async, inject } from '@angular/core/testing';

import { WithPermissionGuard } from './with-permission.guard';

describe('WithPermissionGuard', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [WithPermissionGuard]
    });
  });

  it('should ...', inject([WithPermissionGuard], (guard: WithPermissionGuard) => {
    expect(guard).toBeTruthy();
  }));
});
