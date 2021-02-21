import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ClanDetailComponent } from './clan-detail.component';

describe('ClanDetailComponent', () => {
  let component: ClanDetailComponent;
  let fixture: ComponentFixture<ClanDetailComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ClanDetailComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ClanDetailComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
