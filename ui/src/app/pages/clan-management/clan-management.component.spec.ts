import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ClanManagementComponent } from './clan-management.component';

describe('ClanManagementComponent', () => {
  let component: ClanManagementComponent;
  let fixture: ComponentFixture<ClanManagementComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ClanManagementComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ClanManagementComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
